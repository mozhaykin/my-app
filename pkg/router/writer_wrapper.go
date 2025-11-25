package router

import "net/http"

type Writer struct {
	http.ResponseWriter // эмбединг

	wroteCode bool // был ли записан код ответа
	code      int  // сам код ответа
}

func WriterWrapper(w http.ResponseWriter) *Writer {
	return &Writer{ResponseWriter: w}
}

func (w *Writer) WriteHeader(code int) { // Переопределяем метод WriteHeader, добавляем в него свою логику.
	if !w.wroteCode { // Если код не записан, то:
		w.setCode(code)                    // записываем его в нашу структуру (w.wroteCode = true и w.code = code)
		w.ResponseWriter.WriteHeader(code) // оставляем прежнюю логику
	}
}

func (w *Writer) Write(data []byte) (int, error) { // Переопределяем метод Write, добавляем в него свою логику.
	if !w.wroteCode { // Если код не записан, то:
		w.setCode(http.StatusOK) // записываем его в нашу структуру (w.wroteCode = true и w.code = code)
	}
	// оставляем прежнюю логику
	return w.ResponseWriter.Write(data) //nolint:wrapcheck
}

func (w *Writer) Code() int { // Метод для того чтобы потом можно было достать код
	return w.code
}

func (w *Writer) setCode(code int) {
	w.wroteCode = true
	w.code = code
}
