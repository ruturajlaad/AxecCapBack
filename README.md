🧠 AxeCap — CFD Trading Platform

AxeCap is a backend trading execution engine built using Go (Golang) and Gin and Next.js ,Tailwind CSS, designed to simulate and manage CFD-style leveraged trading with real-time crypto market data (via Binance WebSocket feed).

🚀 Current Features (as of now)
🧩 1. User Authentication System

Secure JWT-based login and registration system.

Routes:

POST /register → create user

POST /login → authenticate and get JWT token

💸 2. Account & Balance Management

Every user has a balance field stored in the database.

When positions close, profit/loss (PnL) is automatically credited/debited from user balance.

You can fetch current account balance via:

GET /user/balance → returns real-time updated balance for logged-in user

📈 3. Order Management

Users can place limit orders through:

POST /order

Orders are stored with “pending” status until the market price reaches the limit price.

⚙️ 4. Real-Time Price Execution Engine

Integrated with Binance WebSocket stream (BTC/USDT or any symbol).

When new price data arrives:

Pending orders are automatically checked.

If price conditions match → order executes instantly.

Execution creates an open position with details like:

Entry Price

Leverage

Contracts

StopLoss / TakeProfit

Timestamp

📊 5. Position Tracking

Positions created on execution are tracked in DB.

You can retrieve open positions using:

GET /positions

Positions automatically close when:

Stop Loss or Take Profit levels are hit (can be extended later).

Or manually through:

POST /close-position

💰 6. Real-Time PnL Computation

When a position closes:

Profit/Loss is calculated based on entry/exit price, leverage, and direction (buy/sell).

Result is reflected immediately in user’s account balance.

🔌 7. Database Layer

Built using clean repository pattern:

Handles CRUD for users, orders, and positions.

Includes database migrations for schema setup (MySQL/Postgres ready).

🧱 8. Modular Clean Architecture

internal/trade/ → handles business logic (service, repo, models, handlers)

internal/api/ → manages routes and HTTP layer

internal/middleware/ → JWT auth middleware

cmd/main.go → entry point wiring everything together

🪙 9. Live Price Integration

Uses your Binance WebSocket feed to stream real-time BTC/USDT price.

This feed powers:

Order triggering

Price display on TradingView charts (can be used both for visualization and logic)

🧭 10. What’s Next (Future Roadmap)

Implement margin requirements & liquidation logic

Add account summary (total equity, used margin, free margin)

Add multi-symbol support (ETH, XRP, etc.)

Integrate with front-end (TradingView chart + order panel)

Add REST/WebSocket server for broadcasting live trades and updates
