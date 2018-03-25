Jixbot Commands
=========================

Jixbot provides a variety of commands, they are documented below. The commands require different statuses in the channel to use:

*B* - Broadcaster or admin   
*M* - Mod or *B*  
*S* - Subscriber or *M*  
*V* - Viewer or *S*  

Some commands require additional settings for the channel. A dashboard for managing these commands and settings is a work in progress.

Management
----------------
* **add** (*M*)
  * `!addcommand !<name> <optional: -cd=seconds> <response>` - Adds a text command (overwrites existing)
  * `!addsubcommand !<name> <optional: -cd=seconds> <response>` - Adds a text commands only usable by subs (overwrites existing)
  * `!addmodcommand !<name> <optional: -cd=seconds> <response>` - Adds a text command only usable by mods (overwrites existing)
    * *-cd= sets a cooldown in seconds on the command*
    * *Example:* `!addcommand !usefulmessage Reminder that subscribing is not necessary to support the streamer!`
    * *Example:* `!addcommand !superspam -cd=30 SPAM SPAM Kappa SPAM`
  * Can use `$u$` to substitute the invoking user, and `$0$`, `$1$`, etc to substitute arguments.
    * *Example:* `!addcommand !userandargs user: $u$, second: $1$, first: $0$`
    * *Example:* `!userandargs Kappa SMOrc SMOrc-> user: jixbot, second: SMOrc SMOrc, first: Kappa`
  * Can use `$url: <url>$` to substitute an api call's response. The above substitutions are applied to the url
    * *Example:* `!addcommand !url $url: http://www.dummyapi.com/text/$0$ $`
    * *Example:* `!url hello -> the text passed into the api was "hello"`
  * Can use `$rand:<min>-<max>$` to substitute a random number in the response
    * *Example:* `!addcommand !rollthedice $u$ rolled a $rand:1-6$`
    * *Example:* `!rollthedice -> jix rolled a 6
* **delete** (*M*)
  * `!deletecommand !<name>` - Deletes a text command
* **modonly** (*M*)
  * `!modonly` - Makes all jixbot commands only usable by mods
  * `!subonly` - Makes all jixbot commands only usable by subs
  * `!commandmode` - Displays the restrictions on jixbot commands

Stats
---------------
* **uptime** (*V*)
  * `!uptime` - Displays how long the stream has been online
* **lines_typed** (*V*)
  * `!linestyped` - Whispers the caller how many lines they've typed in chat
* **time_spent**
  * `!timespent` (*V*) - Whispers the caller how long they've been spent in online chat
  * `!longestwatchers` (*M*) - Displays the users with the highest viewing times
* **money** (*V*)
  * `!cash` - Whispers the caller how much currency they have
  * `!givecash <username> <amount>` - Gives `<amount>` coins to `<username>`. Whispers confirmation to both sender and receiver

Interactive commands
----------------
These commands are meant to create and/or enhance chat interaction.

* **failfish** - When someone mis-capitalizes the FailFish emote, corrects them
* **conversation** - When someone sends a message with "jixbot" in it, jixbot replies in a *usually* friendly manner
* **combo** - Combo amount increases as a particular emote or phrase is said in chat, kept alive by unique chatters who say it. Ends when enough time has passed since the last time combo was kept alive. Default combo emote is PogChamp
* **submessage** - Welcome message is displayed when someone subs. Uses sub emotes, so having jixbot subbed is useful
* **question** - Allows jixbot to answer questions asked in chat
  * `!q @username <answer>` (*M*) - Saves an answer to the question last asked by @username
  * *Jixbot only recognizes questions if they begin with [who,what,when,where,why,how] or end with ?*
  * *Jixbot will answer questions that are similar to what it has stored. However, the more stored questions, the more effective the answering will be*

Misc commands
---------------
* **emotes** (*V*)
  * `!emotes` - Lists emotes in the channel
* **brawl**
  * `!brawl` (*M*) - Starts a brawl. Anyone can join with !pileon, only one person wins
  * `!pileon <optional: bet=amount> <optional: weapon>` (*V*) - Join an active brawl with an optional weapon
    * *bet= adds a bet on yourself winning*
    * *Example:* `!pileon bet=1337 Lethal Weapon`
  * `!brawlstats <season or "all">` (*V*) - Display top brawl winners of season or overall
  * `!brawlwins <optional: season>` (*V*) - Show how many brawl wins the user has in a particular season (defaults to current season)
  * `!newbrawlseason` (*B*) - Start a new brawl season
* **quotes**
  * `!addquote` (*M*) - Add a text quote with the remaining text in the message
  * `!quote` (*V*) - Get a quote.
    * `!quote list` will dump all quotes into pastebin and post the link in chat.
    * If a quote number is supplied, that quote will be retrieved if it exists
    * If some text is supplied, a quote will be retrieved that contains that text
    * If no argument is supplied, a random quote will be retrieved
  * `!deletequote <quote number>` (*M*) - Delete a quote
  * `!addclip` (*M*) - Add a clip with the remaining text in the message
  * `!clip` (*V*) - Get a clip.
    * `!clip list` will dump all clips into pastebin and post the link in chat.
    * If a clip number is supplied, that clip will be retrieved if it exists
    * If some text is supplied, a clip will be retrieved that contains that text
    * If no argument is supplied, a random clip will be retrieved
  * `!deletclip <clip number>` (*M*) - Delete a clip



