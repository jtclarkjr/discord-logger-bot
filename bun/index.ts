import 'dotenv/config'
import { Hono } from 'hono'
import { Client, GatewayIntentBits, Message, TextChannel } from 'discord.js'
import fs from 'fs'
import path from 'path'

let client: Client | null = null

const LOG_FILE = path.resolve(process.cwd(), 'mod_logs.txt')

// Store recent messages to log deletion details
const messageCache = new Map<
  string,
  { content: string; author: string; channel: string; timestamp: Date }
>()
const CACHE_DURATION = 24 * 60 * 60 * 1000 // 24 hours

/**
 * Logs a message to a local file for moderation.
 * Includes timestamp, guild, channel, author, and content.
 */
function logMessage(message: Message) {
  // Only log regular user messages from guild text channels
  if (!message.guild || message.author.bot) return

  // Cache message for potential deletion tracking
  messageCache.set(message.id, {
    content: message.content,
    author: `${message.author.tag}`,
    channel: message.channel instanceof TextChannel ? message.channel.name : 'Unknown',
    timestamp: new Date()
  })

  const logEntry = [
    `[${new Date().toISOString()}]`,
    `[MESSAGE]`,
    `Channel: ${message.channel instanceof TextChannel ? message.channel.name : 'Unknown'}`,
    `Author: ${message.author.tag}`,
    `Content: ${message.content.replace(/\n/g, '\\n')}`,
    '\n'
  ].join(' | ')

  fs.appendFileSync(LOG_FILE, logEntry)
  console.log(logEntry)
}

/**
 * Logs a deleted message to the file.
 */
function logDeletedMessage(messageId: string, channelName: string) {
  const cachedMessage = messageCache.get(messageId)

  const logEntry = [
    `[${new Date().toISOString()}]`,
    `[DELETED]`,
    `Channel: ${channelName}`,
    cachedMessage ? `Author: ${cachedMessage.author}` : 'Author: Unknown',
    cachedMessage
      ? `Content: ${cachedMessage.content.replace(/\n/g, '\\n')}`
      : 'Content: Unknown (not cached)',
    '\n'
  ].join(' | ')

  fs.appendFileSync(LOG_FILE, logEntry)
  console.log(logEntry)

  // Remove from cache after logging
  messageCache.delete(messageId)
}

/**
 * Cleans up old cached messages to prevent memory leaks.
 */
function cleanupMessageCache() {
  const now = Date.now()
  for (const [messageId, data] of messageCache.entries()) {
    if (now - data.timestamp.getTime() > CACHE_DURATION) {
      messageCache.delete(messageId)
    }
  }
}

/**
 * Starts the Discord bot.
 */
async function startBot(): Promise<string> {
  if (client) return 'Bot is already running.'

  client = new Client({
    intents: [
      GatewayIntentBits.Guilds,
      GatewayIntentBits.GuildMessages,
      GatewayIntentBits.MessageContent
    ]
  })

  // Log every new message (except from bots)
  client.on('messageCreate', (message: Message) => {
    logMessage(message)
  })

  // Log deleted messages
  client.on('messageDelete', (message) => {
    // Only log deletions from guild text channels
    if (!message.guild) return

    const channelName = message.channel instanceof TextChannel ? message.channel.name : 'Unknown'
    logDeletedMessage(message.id, channelName)
  })

  // Clean up message cache every hour
  setInterval(cleanupMessageCache, 60 * 60 * 1000)

  client.once('ready', () => {
    console.log(`Logged in as ${client?.user?.tag}!`)
  })

  try {
    await client.login(process.env.DISCORD_BOT_TOKEN)
    return 'Bot started successfully.'
  } catch (error) {
    console.error('Error logging in:', error)
    client = null
    return 'Failed to start bot.'
  }
}

async function stopBot(): Promise<string> {
  if (!client) return 'Bot is not running.'
  try {
    await client.destroy()
    client = null
    return 'Bot stopped successfully.'
  } catch (error) {
    console.error('Error stopping bot:', error)
    return 'Failed to stop bot.'
  }
}

/**
 * Hono REST endpoints
 */

const app = new Hono()
const port = process.env.PORT || 8080

app.post('/bot/mod-logger/on', async (context) => {
  const result = await startBot()
  return context.text(result, 200)
})

app.post('/bot/mod-logger/off', async (context) => {
  const result = await stopBot()
  return context.text(result, 200)
})

// For Bun/Node serverless:
export default {
  port,
  fetch: app.fetch
}
