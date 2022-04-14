# GO serv with request log

1. Start postgres container:
`docker-compose up --build`
2. Create DB (migrations)
`go run ./db`
3. Start server: `go run .`
4. To PUT example
`
POST http://localhost:4000/payment
Content-Type: application/json
{
"resep_id": 1,
"val": 600
}`
5. To TRANSFER example
`
POST http://localhost:4000/transfer
Content-Type: application/json
{
"sender_id": 1,
"resep_id": 2,
"val": 1600
}
`