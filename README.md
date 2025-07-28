# Discord Logger Bot

A Discord bot designed for server moderation and logging purposes. This bot helps moderators track and monitor server activity for maintaining a safe community environment.

## ⚠️ Important Notice

**This bot is intended for moderation purposes only.** It should be used by server administrators and moderators to maintain server safety and enforce community guidelines. Please ensure you comply with Discord's Terms of Service and Privacy Policy when using this bot.

## Discord Bot Setup

### 1. Create a Discord Application
1. Go to the [Discord Developer Portal](https://discord.com/developers/applications)
2. Click "New Application" and give it a name
3. Navigate to the "Bot" section in the left sidebar
4. Click "Add Bot" to create a bot user
5. Copy the bot token (you'll need this for configuration)

### 2. Set Bot Permissions
1. In the "Bot" section, scroll down to "Privileged Gateway Intents"
2. Enable the following intents based on your logging needs:
   - Message Content Intent (for message logging)
   - Server Members Intent (for member activity logging)
3. Go to the "OAuth2" > "URL Generator" section
4. Select "bot" scope and required permissions:
   - Read Messages/View Channels
   - Send Messages
   - Read Message History
   - Manage Messages (if needed for moderation)

### 3. Invite Bot to Server
1. Copy the generated OAuth2 URL
2. Open it in your browser and select your server
3. Authorize the bot with the selected permissions

## Configuration

### 1. Configuration File
Update the configuration file with your server-specific settings:
- Log channels for different types of events
- Moderation settings
- Logging preferences

## Features

- Message logging and moderation
- User activity monitoring
- Server event tracking
- Moderation action logging

## API Endpoints

### Mod Logger Control

You can control the mod logger functionality using HTTP endpoints:

#### Enable Mod Logger
- **Production**: `POST /bot/mod-logger/on`
- **Local Development**: `POST localhost:8080/bot/mod-logger/on`

#### Disable Mod Logger
- **Production**: `POST /bot/mod-logger/off`
- **Local Development**: `POST localhost:8080/bot/mod-logger/off`

These endpoints allow you to dynamically enable or disable the moderation logging features without restarting the bot.

## Moderation Usage Only

**This bot is designed exclusively for server moderation purposes and should only be used by:**
- Server administrators
- Appointed moderators
- Users with explicit permission for moderation activities

**Appropriate use cases include:**
- Monitoring for rule violations
- Tracking spam or abuse
- Maintaining server safety
- Enforcing community guidelines

**This bot should NOT be used for:**
- Personal surveillance
- Harassment or stalking
- Privacy violations
- Non-moderation related monitoring

## Privacy and Compliance

### Discord Privacy Policy Compliance

This bot operates in accordance with Discord's Terms of Service and Privacy Policy. Users should be aware that:

1. **Data Collection**: This bot may collect and log messages, user activity, and server events for moderation purposes only.

2. **Data Storage**: Logged data is stored securely and used exclusively for server moderation and safety purposes. The bot creates a `mod_logs.txt` file to store moderation events and also outputs logs to the terminal for real-time monitoring.

3. **Data Retention**: Log data should be retained only as long as necessary for moderation purposes and in compliance with Discord's policies.

4. **User Rights**: Users have the right to request information about data collection and may request data deletion where applicable.

5. **Transparency**: Server members should be informed about the presence and purpose of logging bots through server rules or announcements.

### Data Storage Locations

- **File Logging**: Moderation events are stored in `mod_logs.txt` for persistent record-keeping
- **Terminal Logging**: Real-time activity is displayed in the console/terminal for immediate monitoring
- **Secure Storage**: All log files should be stored securely and access should be restricted to authorized moderators only

### Compliance Requirements

- Ensure your server's privacy policy mentions the use of logging bots
- Inform users about data collection practices
- Use collected data only for legitimate moderation purposes
- Regularly review and clean up stored logs (including `mod_logs.txt`)
- Secure access to log files and terminal output
- Respect user privacy and Discord's community guidelines

## Legal Disclaimer

By using this bot, you agree to:
- Use it only for legitimate server moderation purposes
- Comply with Discord's Terms of Service and Privacy Policy
- Respect user privacy and applicable data protection laws
- Take responsibility for proper configuration and usage

**Repository Disclaimer**: This repository and its maintainers do not take responsibility for any actions, consequences, or damages resulting from the use, misuse, or implementation of this bot. Users are solely responsible for ensuring compliance with all applicable laws, platform policies, and ethical guidelines when deploying and operating.


---

**Remember**: Always prioritize user privacy and community safety when using moderation tools.
