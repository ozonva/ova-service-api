package api_test

import (
	"context"
	"fmt"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/google/uuid"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/ozonva/ova-service-api/internal/api"
	"github.com/ozonva/ova-service-api/internal/mocks"
	"github.com/ozonva/ova-service-api/internal/models"

	pb "github.com/ozonva/ova-service-api/pkg/ova-service-api"
)

var _ = Describe("Api", func() {
	var (
		ctx         context.Context
		ctrl        *gomock.Controller
		repoMock    *mocks.MockRepo
		flusherMock *mocks.MockFlusher
		saverMock   *mocks.MockSaver

		carServiceID string
		carService   models.Service

		validMultiCreateRequest []*pb.CreateServiceV1Request
	)

	BeforeEach(func() {
		ctx = context.Background()
		ctrl = gomock.NewController(GinkgoT())
		repoMock = mocks.NewMockRepo(ctrl)
		flusherMock = mocks.NewMockFlusher(ctrl)
		saverMock = mocks.NewMockSaver(ctrl)

		carServiceID = "d6fa505c-6072-4a45-bdae-86e6b13d7342"
		carService = models.Service{
			ID:          uuid.MustParse(carServiceID),
			UserID:      1,
			ServiceName: "Car service",
		}

		validMultiCreateRequest = []*pb.CreateServiceV1Request{
			{UserId: 1, ServiceName: "Panzer service"},
			{UserId: 1, ServiceName: "Yacht service"},
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Describe("API endpoint handlers to save and load services using repo", func() {
		Context("on calling Create endpoint", func() {
			When("request body is empty", func() {
				It("should return InvalidArgument error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().AddServices(gomock.Any()).Times(0)

					_, err := server.CreateServiceV1(ctx, nil)

					Expect(err).Should(HaveOccurred())
				})
			})

			When("request body contains illegal service data", func() {
				It("should return Internal error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().AddServices(gomock.Any()).Times(0)

					_, err := server.CreateServiceV1(ctx, &pb.CreateServiceV1Request{UserId: 0})

					Expect(err).Should(HaveOccurred())
				})
			})

			When("repo return error", func() {
				It("should return Internal error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().AddServices(gomock.Any()).
						Return(fmt.Errorf("repo err")).Times(1)

					_, err := server.CreateServiceV1(ctx, &pb.CreateServiceV1Request{UserId: 1})

					Expect(err).Should(HaveOccurred())
				})
			})

			When("valid request", func() {
				It("should return serviceID", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().AddServices(gomock.Any()).Times(1)

					res, err := server.CreateServiceV1(ctx, &pb.CreateServiceV1Request{UserId: 1})

					Expect(err).ShouldNot(HaveOccurred())
					Expect(res.ServiceId).ShouldNot(BeEmpty())
				})
			})
		})

		Context("on calling Describe endpoint", func() {
			When("request body is empty", func() {
				It("should return InvalidArgument error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().DescribeService(gomock.Any()).Times(0)

					_, err := server.DescribeServiceV1(ctx, nil)

					Expect(err).Should(HaveOccurred())
				})
			})

			When("can't parse serviceID", func() {
				It("should return InvalidArgument error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().DescribeService(gomock.Any()).Times(0)

					_, err := server.DescribeServiceV1(ctx, &pb.DescribeServiceV1Request{ServiceId: "bad uuid"})

					Expect(err).Should(HaveOccurred())
				})
			})

			When("service not found", func() {
				It("should return NotFound error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().DescribeService(gomock.Any()).
						Return(nil, fmt.Errorf("not found")).Times(1)

					_, err := server.DescribeServiceV1(ctx, &pb.DescribeServiceV1Request{ServiceId: carServiceID})

					Expect(err).Should(HaveOccurred())
				})
			})

			When("can't map service model to response", func() {
				It("should return Internal error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().DescribeService(gomock.Any()).Return(nil, nil).Times(1)

					_, err := server.DescribeServiceV1(ctx, &pb.DescribeServiceV1Request{ServiceId: carServiceID})

					Expect(err).Should(HaveOccurred())
				})
			})

			When("valid request", func() {
				It("should return service", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().DescribeService(gomock.Any()).Return(&carService, nil).Times(1)

					res, err := server.DescribeServiceV1(ctx, &pb.DescribeServiceV1Request{ServiceId: carServiceID})

					Expect(err).ShouldNot(HaveOccurred())
					Expect(res.ServiceId).Should(BeEquivalentTo(carServiceID))
				})
			})
		})

		Context("on calling List endpoint", func() {
			When("error occurs in repo", func() {
				It("should return Internal error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().ListServices(gomock.Any(), gomock.Any()).
						Return(nil, fmt.Errorf("repo error")).Times(1)

					_, err := server.ListServicesV1(ctx, &empty.Empty{})

					Expect(err).Should(HaveOccurred())
				})
			})

			When("valid request", func() {
				It("should return list of services", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().ListServices(gomock.Any(), gomock.Any()).
						Return([]models.Service{carService, carService}, nil).Times(1)

					res, err := server.ListServicesV1(ctx, &empty.Empty{})

					Expect(err).ShouldNot(HaveOccurred())
					Expect(len(res.ServiceShortInfo)).Should(BeEquivalentTo(2))
					Expect(res.ServiceShortInfo[0].ServiceId).Should(BeEquivalentTo(carServiceID))
					Expect(res.ServiceShortInfo[1].ServiceId).Should(BeEquivalentTo(carServiceID))
				})
			})
		})

		Context("on calling Remove endpoint", func() {
			When("request body is empty", func() {
				It("should return InvalidArgument error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().RemoveService(gomock.Any()).Times(0)

					_, err := server.RemoveServiceV1(ctx, nil)

					Expect(err).Should(HaveOccurred())
				})
			})

			When("can't parse serviceID", func() {
				It("should return InvalidArgument error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().RemoveService(gomock.Any()).Times(0)

					_, err := server.RemoveServiceV1(ctx, &pb.RemoveServiceV1Request{ServiceId: "bad uuid"})

					Expect(err).Should(HaveOccurred())
				})
			})

			When("service not found", func() {
				It("should return NotFound error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().RemoveService(gomock.Any()).
						Return(fmt.Errorf("not found")).Times(1)

					_, err := server.RemoveServiceV1(ctx, &pb.RemoveServiceV1Request{ServiceId: carServiceID})

					Expect(err).Should(HaveOccurred())
				})
			})

			When("valid request", func() {
				It("should return empty result after removing", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().RemoveService(gomock.Any()).
						Return(nil).Times(1)

					res, err := server.RemoveServiceV1(ctx, &pb.RemoveServiceV1Request{ServiceId: carServiceID})

					Expect(err).ShouldNot(HaveOccurred())
					Expect(res).Should(BeEquivalentTo(&empty.Empty{}))
				})
			})
		})

		Context("on calling MultiCreate endpoint", func() {
			When("request body is empty", func() {
				It("should return InvalidArgument error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					flusherMock.EXPECT().Flush(gomock.Any()).Times(0)

					_, err := server.MultiCreateServiceV1(ctx, nil)

					Expect(err).Should(HaveOccurred())
				})
			})

			When("request body contains list with invalid objects", func() {
				It("should return Argument error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					flusherMock.EXPECT().Flush(gomock.Any()).Times(0)
					req := &pb.MultiCreateServiceV1Request{CreateService: []*pb.CreateServiceV1Request{nil}}

					_, err := server.MultiCreateServiceV1(ctx, req)

					Expect(err).Should(HaveOccurred())
				})
			})

			When("can't flush all services to repo", func() {
				It("should return Internal error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					flusherMock.EXPECT().Flush(gomock.Any()).
						Return([]models.Service{carService}).Times(1)
					req := &pb.MultiCreateServiceV1Request{CreateService: validMultiCreateRequest}

					_, err := server.MultiCreateServiceV1(ctx, req)

					Expect(err).Should(HaveOccurred())
				})
			})

			When("valid request", func() {
				It("should return slice of serviceID", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					flusherMock.EXPECT().Flush(gomock.Any()).Return(nil).Times(1)
					req := &pb.MultiCreateServiceV1Request{CreateService: validMultiCreateRequest}

					res, err := server.MultiCreateServiceV1(ctx, req)

					Expect(err).ShouldNot(HaveOccurred())
					Expect(res).Should(BeEquivalentTo(&empty.Empty{}))
				})
			})
		})

		Context("on calling Update endpoint", func() {
			When("request body is empty", func() {
				It("should return InvalidArgument error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().UpdateService(gomock.Any()).Times(0)

					_, err := server.UpdateServiceV1(ctx, nil)

					Expect(err).Should(HaveOccurred())
				})
			})

			When("can't parse serviceID", func() {
				It("should return InvalidArgument error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().UpdateService(gomock.Any()).Times(0)

					_, err := server.UpdateServiceV1(ctx, &pb.UpdateServiceV1Request{ServiceId: "bad uuid"})

					Expect(err).Should(HaveOccurred())
				})
			})

			When("request body contains illegal service data", func() {
				It("should return Internal error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().UpdateService(gomock.Any()).Times(0)

					_, err := server.UpdateServiceV1(ctx, &pb.UpdateServiceV1Request{ServiceId: carServiceID, UserId: 0})

					Expect(err).Should(HaveOccurred())
				})
			})

			When("repo return error", func() {
				It("should return Internal error", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().UpdateService(gomock.Any()).
						Return(fmt.Errorf("repo err")).Times(1)

					_, err := server.UpdateServiceV1(ctx, &pb.UpdateServiceV1Request{ServiceId: carServiceID, UserId: 1})

					Expect(err).Should(HaveOccurred())
				})
			})

			When("valid request", func() {
				It("should update service", func() {
					server := api.NewGrpcApiServer(repoMock, saverMock, flusherMock)
					repoMock.EXPECT().UpdateService(gomock.Any()).Times(1)

					res, err := server.UpdateServiceV1(ctx, &pb.UpdateServiceV1Request{ServiceId: carServiceID, UserId: 1})

					Expect(err).ShouldNot(HaveOccurred())
					Expect(res).Should(BeEquivalentTo(&empty.Empty{}))
				})
			})
		})
	})
})
