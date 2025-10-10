package logger

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"

	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/internal/dto/baggage"
	"gitlab.golang-school.ru/potok-1/amozhaykin/my-app/pkg/router"
)

func Middleware(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		bag := &baggage.Baggage{}
		// содаем дочерний контекст с value. В качестве value пустая структура bag, которую мы запоним вызывая контекст
		// в нужном месте. Например в render.error мы положим туда ошибку
		ctx := baggage.WithContext(r.Context(), bag)

		ww := router.WriterWrapper(w)          // создаем обертку над стадартным http.ResponseWriter
		next.ServeHTTP(ww, r.WithContext(ctx)) // возвращаем ServeHTTP с обернутым http.ResponseWriter

		// event (событие) это запись лога
		// по умолчанию создается лог INFO
		event := log.Info()

		if bag.Err != nil {
			event = log.Error().Err(bag.Err) // если есть ошибка, то лог меняется на ERROR и ошибка записывается в него
		}

		event. // записываем в лог остальные данные, которые нам нужны
			Str("profile_id", bag.ProfileID).
			Str("proto", "http").
			Int("code", ww.Code()).
			Str("method", fmt.Sprintf("%s %s", r.Method, router.ExtractPath(ctx))).
			Send()
	}

	return http.HandlerFunc(fn)
}
