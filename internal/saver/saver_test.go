package saver_test

import (
	"golang.org/x/net/context"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozonva/ova-service-api/internal/mocks"
	"github.com/ozonva/ova-service-api/internal/models"

	saver_ "github.com/ozonva/ova-service-api/internal/saver"
)

const (
	shortTimeout = 2 * time.Second
	longTimeout  = 10 * time.Second
	finalTimeout = 1 * time.Second
)

var _ = Describe("Saver", func() {
	var (
		ctrl          *gomock.Controller
		flusherMock   *mocks.MockFlusher
		carService    models.Service
		panzerService models.Service
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		flusherMock = mocks.NewMockFlusher(ctrl)

		carService = models.Service{ID: uuid.New(), ServiceName: "Car service"}
		panzerService = models.Service{ID: uuid.New(), ServiceName: "Panzer service"}
	})

	AfterEach(func() {
		// Required to be sure that flush goroutine has a chance to run
		time.Sleep(finalTimeout)
		ctrl.Finish()
	})

	Describe("Save services using Saver", func() {
		Context("on Save service", func() {
			When("local storage have enough space to store the service", func() {
				It("should not call flusher immediately", func() {
					saver := saver_.New(1, longTimeout, flusherMock)
					saver.Init()

					flusherMock.EXPECT().Flush(gomock.Any(), gomock.Any()).Times(0)

					_ = saver.Save(carService)
				})
			})

			When("local storage is full", func() {
				It("should return error", func() {
					saver := saver_.New(1, longTimeout, flusherMock)
					saver.Init()

					flusherMock.EXPECT().Flush(gomock.Any(), gomock.Eq(gomock.Any())).Times(0)

					_ = saver.Save(carService)
					Expect(saver.Save(panzerService)).Should(HaveOccurred())
				})
			})
		})

		Context("on timeout expiration", func() {
			It("should flush when timeout expired", func() {
				saver := saver_.New(1, shortTimeout, flusherMock)
				saver.Init()

				flusherMock.EXPECT().Flush(context.Background(), gomock.Eq([]models.Service{carService})).Times(1)

				_ = saver.Save(carService)
				time.Sleep(shortTimeout)
			})
		})

		Context("on Close saver", func() {
			It("should flush", func() {
				saver := saver_.New(1, longTimeout, flusherMock)
				saver.Init()

				flusherMock.EXPECT().Flush(context.Background(), gomock.Eq([]models.Service{carService})).Times(1)

				_ = saver.Save(carService)
				saver.Close()
			})
		})
	})
})
