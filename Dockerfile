#build stage
from golang:1.20.3-alpine  AS builder
# run apk add --no.cache git upx
run apk add --no-cache git upx
workdir /app

copy ["go.mod","go.sum","./"]

copy . .


run go get -d -v .

run go build \
  -ldflags="-s -w" \
  -o app ./ 

# run go build ./cmd/novelas/main.go \
#     -ldflags="-s -w" \
#     -o app -v -

run upx app



#final stage

from alpine 
label Name=Scraping

run apk update
run apk --no-cache add ca-certificates

workdir /app

copy --from=builder /app .

ENTRYPOINT ["./app"]
     