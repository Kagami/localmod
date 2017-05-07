all: build serve

config:
	cp localmod.yaml.example localmod.yaml

serve:
	./localmod -config localmod.yaml serve

build: clean
	go build -ldflags="-s -w"

clean:
	rm -f localmod
