# CopaDonation

https://santimoncada.github.io/CopaDonation/

## How to run

```
npm install

go mod tidy

npm run build

go run src/*.go

```

## how to develop

```
nodemon --exec go run src/*.go --signal SIGTERM

npm run dev
```

## todo list

- [x] Donate button
- [x] List Stripe donations
- [x] Webhooks
- [x] JS bundling
- [ ] Make concurrent calls to stripe
- [ ] Cache templates on startup and webhook call
- [ ] HTTP streaming
- [ ] Design
- [ ] CSS
