package errorhandler

import (
	"log"
	"net/http"
)

// 用这个自定义类包装处理器方法，被包装的方法必须返回一个error，这个自定义类作为拦截器统一将error转为http响应
// 以此实现类似SpringBoot的异常统一处理效果。
type ErrorToHttpResponse func(http.ResponseWriter, *http.Request) error

// 实现 http.Handler 接口
func (fn ErrorToHttpResponse) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn(w, r); err != nil {
		log.Printf("ErrorToHttpResponse 统一处理错误: %v", err)

		// 针对 BusinessError 处理一下 httpStatus
		httpStatus := http.StatusInternalServerError
		if businessError, ok := err.(*BusinessError); ok && businessError.HttpCode != 0 {
			httpStatus = businessError.HttpCode
		}
		http.Error(w, err.Error(), httpStatus)
	}
}