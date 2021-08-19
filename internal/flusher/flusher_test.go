package flusher_test

import (
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozonva/ova-service-api/internal/mocks"
	"github.com/ozonva/ova-service-api/internal/models"

	flusher_ "github.com/ozonva/ova-service-api/internal/flusher"
)

var _ = Describe("Flusher", func() {
	var (
		ctrl     *gomock.Controller
		repoMock *mocks.MockRepo
		services []models.Service
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		repoMock = mocks.NewMockRepo(ctrl)

		services = []models.Service{
			{ID: uuid.New(), ServiceName: "Car service"},
			{ID: uuid.New(), ServiceName: "Bicycle service"},
			{ID: uuid.New(), ServiceName: "Panzer service"},
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("Flush services to the repo", func() {
		Context("Services can be properly split to batches", func() {
			When("Repo able to save every single batch", func() {
				It("flusher.Flush should return nil", func() {
					flusher := flusher_.New(2, repoMock)

					gomock.InOrder(
						repoMock.EXPECT().AddServices(gomock.Eq(services[0:2])).Return(nil),
						repoMock.EXPECT().AddServices(gomock.Eq(services[2:])).Return(nil),
					)
					Expect(flusher.Flush(services)).To(BeNil())
				})
			})

			When("Repo able to save first batch, but failed on second", func() {
				It("flusher.Flush should return second batch", func() {
					flusher := flusher_.New(2, repoMock)

					gomock.InOrder(
						repoMock.EXPECT().AddServices(gomock.Eq(services[0:2])).Return(nil),
						repoMock.EXPECT().AddServices(gomock.Eq(services[2:])).Return(fmt.Errorf("connection failed")),
					)
					Expect(flusher.Flush(services)).To(BeEquivalentTo(services[2:]))
				})
			})

			When("Repo able to save second batch, but failed on first", func() {
				It("flusher.Flush should return first batch", func() {
					flusher := flusher_.New(2, repoMock)

					gomock.InOrder(
						repoMock.EXPECT().AddServices(gomock.Eq(services[0:2])).Return(fmt.Errorf("connection failed")),
						repoMock.EXPECT().AddServices(gomock.Eq(services[2:])).Return(nil),
					)
					Expect(flusher.Flush(services)).To(BeEquivalentTo(services[0:2]))
				})
			})
		})

		Context("Batch size is zero value", func() {
			It("flusher.Flush should return entire services slice", func() {
				flusher := flusher_.New(0, repoMock)
				Expect(flusher.Flush(services)).To(BeEquivalentTo(services))
			})
		})
	})
})
