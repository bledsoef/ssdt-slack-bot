# Student Software Development Team Slack Bot
## Purpose
The purpose of this application is to make the Berea student programmer's time more efficient. Programmer typically have 1.5 - 2.5 hour shifts which is not much time. It feels like even less time when you need to spend time combing back through your PRs to see which of yours need changes and which of them still require review. It can especially be slow for peer reviewers who have to spend much of their shift looking through many PRs that they may not need to look at. This bot will send a message every morning telling you what PRs need review and what PRs need changes. It will tag all of the people assigned to each PR and it will monitor who is currenly reviewing PRs and tag them as well. So now, instead of having to check multiple sources for what PRs require review, who is on PR duty, and additionally who is assigned to what cleaning tasks, one check in to Slack will tell you all the information you need to know to help you get started ASAP.

## Usage
This application has two major functionalities, PR reviewing and chores.

### PR Reviewing
Automatically, each day (Monday-Friday) the bot will send a message using Cron to the channel detailing the current PR reviewers, the PRs which need review, and those that require changes.

If you can't find that message or want to see it again you can always enter `/prs` which will send you an ephemeral message (one only you can see) that will contain all of the information that would have been sent by the bot each morning.

### Cleaning
Each Friday at 8am the bot will send a message using Cron with the current week's cleaning schedule and it will tag all of the users associated with each task. 

Like with PRs, if you need to do your chores before or after Friday or perhaps lost the message somehow `/cleaning` is a good way to get all the information that would have been sent on Friday at 8am.

## References
Many external sources were used in the making of this application:
- https://www.bacancytechnology.com/blog/develop-slack-bot-using-golang - helped me initialize my application and start using basic event handling.
- https://github.com/xNok/slack-go-demo-socketmode - allowed me to implement slash commands and a more concise version of Socketmode into my code. 
- https://chat.openai.com/chat - helped me learn things about the Golang-Slack documentation I couldn't find on my own.
- Brian Ramsay - helped me set up Cron and integrate code into SSDT Slack.
