TARGETS := \
	   regexp-match \
	   regexp-quote \
	   regexp-replace \
	   regexp-walk

all: $(TARGETS)

install:
	install -Dm755 -t "$(DESTDIR)/$(PREFIX)/bin" $(TARGETS)

clean:
	rm -f $(TARGETS)

$(TARGETS): %: %.go
	go build -trimpath -o $@ $<
