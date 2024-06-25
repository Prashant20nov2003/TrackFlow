include .env
export $(shell sed 's/=.*//' .env)

server:
	@cd cmd/gotracker && go build && ./gotracker -ip 123.123.123.123

dashboard:
	@cd cmd/dashboard && \
	go build -o localdash && \
	./localdash -site 1 -start 20231101 -end 20231130

build:
	@sed -i 's/http:\/\/localhost:9876/https:\/\/analytics.focuscentric.com/g' src/track.ts
	@npm run build
	@sed -i 's/https:\/\/analytics.focuscentric.com/http:\/\/localhost:9876/g' src/track.ts
	@cd cmd/gotracker && go build -o gotracker
	@echo "Build successful"