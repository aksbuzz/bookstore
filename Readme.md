Contains 3 services

1. Books
2. Cart
3. Orders

## Books

- List available books
- Get a book by book_id
- Get top 10 best selling books

## Cart

- List items in cart
- CRUD operations on cart

## Orders

- List placed orders
- Get order by order_id
- Checkout

### NATS

- NATS request-reply pattern is used for Checkout.
  - Checkout will create a new order and publish an event __order.placed__
  - Cart service will subscribe to this event and remove cart items

### Redis

- Redis Sorted-Sets is used to maintain best sellers

### TODO

- Use NATS Jetstream to update the best sellers

