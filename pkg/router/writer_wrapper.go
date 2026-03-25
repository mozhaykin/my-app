package router

import "net/http"

type Writer struct {
	http.ResponseWriter

	wroteCode bool // был ли записан код ответа
	code      int
}

func WriterWrapper(w http.ResponseWriter) *Writer {
	return &Writer{ResponseWriter: w}
}

func (w *Writer) WriteHeader(code int) { // Переопределяем метод WriteHeader, добавляем в него свою логику.
	if !w.wroteCode {
		w.setCode(code)
		w.ResponseWriter.WriteHeader(code)
	}
}

func (w *Writer) Write(data []byte) (int, error) { // Переопределяем метод Write, добавляем в него свою логику.
	if !w.wroteCode {
		w.setCode(http.StatusOK)
	}

	return w.ResponseWriter.Write(data) //nolint:wrapcheck
}

func (w *Writer) Code() int { // Метод для того чтобы потом можно было достать код
	return w.code
}

func (w *Writer) setCode(code int) {
	w.wroteCode = true
	w.code = code
}
