package middleware

import (
	"my_ecommerce_system/pkg/errorhandler"
	"net/http"
)

func ErrorToHttpHandlingMiddleware(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request){
		defer func() {
			if err := recover(); err != nil{
				errorhandler.HandleErrorToHttpResponse(w, err)
			}
		}()
		next.ServeHTTP(w, r)
	});
}