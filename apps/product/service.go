package product

import (
	"context"
	"fmt"

	"github.com/Tooomat/go-online-shop/infrastructure/response"
	logging "github.com/Tooomat/go-online-shop/internal/log"
)

type productService struct {
	repo ProductRepository
}

func newProductService(repo ProductRepository) productService {
	return productService{
		repo: repo,
	}
}

func (ps productService) CreateProductService(c context.Context, req CreateProductRequestPayload) (err error) {
	productEntity := NewProductFromCreateProductRequest(req)

	if err = productEntity.ProductIsValid(); err != nil {
		//error location from productService to middlerware
		errMsgLoc := fmt.Errorf("[CreateProductService(), ProductIsValid()] err => %v", err)
		logging.ErrorLoggingFromContext(c, errMsgLoc)
		return 
	}

	err = ps.repo.CreateProductRepository(c, productEntity)

	return
}

func (ps productService) ListProductsService(c context.Context, req ListProductRequestPayload) (productEntity []ProductEntity, err error) {
	pagination := req.GenerateDefaultValueRequest()

	productEntity, err = ps.repo.GetAllProductsReporsitoryWithPaginationCursor(c, pagination)
	if err != nil {
		if err == response.ErrNotFound {
			return []ProductEntity{}, nil
		}
		return
	}

	if len(productEntity) == 0 {
		return []ProductEntity{}, nil
	}
	return
}

func (ps productService) ProductDetailService(c context.Context, sku string) (productEntity ProductEntity, err error) {
	productEntity, err = ps.repo.GetProductBySKURepository(c, sku)
	if err != nil {
		return
	}
	return
}
