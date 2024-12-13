# Telegram Bot with Cloudflare AI Integration

This is a Go-based Telegram bot that uses Cloudflare's AI to respond to user inputs in a casual and relatable manner. The bot maintains user-specific chat histories for personalized interactions.

## Features
- Personalized, relaxed responses based on user input.
- Supports private and group chats.
- Retrieves responses from Cloudflare AI endpoint.
- Stores user conversation history for continuous interaction.

## Prerequisites

Before running the bot, make sure you have the following:

- Go 1.18 or higher installed.
- A Telegram bot token (you can get it by talking to [BotFather](https://core.telegram.org/bots#botfather)).
- Cloudflare API token and account ID to interact with their AI.

## Setup Instructions

### Step 1: Install Dependencies

Ensure Go is installed, then run the following command to install necessary dependencies:

```bash
go mod tidy
```

This will fetch all the required packages specified in the project.

### Step 2: Configure Environment Variables

Create a `.env` file in the root directory of the project and add the following variables:

```env
BOT_TOKEN=your_telegram_bot_token
CF_ID=your_cloudflare_account_id
CF_TOKEN=your_cloudflare_api_token
```

- Replace `your_telegram_bot_token` with the Telegram bot token you got from BotFather.
- Replace `your_cloudflare_account_id` and `your_cloudflare_api_token` with the corresponding Cloudflare values.

### Step 3: Run the Bot

After setting up your environment variables, run the bot using:

```bash
go run main.go
```
Or
```bash
./deployment.sh
```