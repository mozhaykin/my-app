package logger

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"github.com/mozhaykin/my-app/internal/dto/baggage"
	"github.com/mozhaykin/my-app/pkg/router"
)

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		bag := &baggage.Baggage{}
		// содаем дочерний контекст с value. В качестве value пустая структура bag, которую мы заполним вызывая контекст
		// в нужном месте. Например в render.error мы положим туда ошибку
		ctx := baggage.WithContext(r.Context(), bag)

		// можно выполнить какую то работу, до вызова ручки

		// создаем обертку над стадартным http.ResponseWriter, для того чтобы потом
		// достать code ответа и записать в log.Info (другой возможности достать код ответа нет)
		ww := router.WriterWrapper(w)

		// дераем нашу ручку (вызываем handler с обернутым http.ResponseWriter)
		next.ServeHTTP(ww, r.WithContext(ctx))

		// выполняется работа после выполнения нашей ручки

		// event это запись лога, по умолчанию создается лог INFO
		event := log.Info()

		if bag.Err != nil {
			event = log.Error().Err(bag.Err) // если есть ошибка, то лог меняется на ERROR и ошибка записывается в него
		}

		event. // записываем в лог остальные данные, которые могут пригодиться для дебагинга
			Str("profile_id", bag.ProfileID).
			Str("proto", "http"). // если у нас есть другие протоколы
			Int("code", ww.Code()).
			Str("method", fmt.Sprintf("%s %s", r.Method, router.ExtractPath(ctx))).
			Send()
	}

	return http.HandlerFunc(fn)
}
