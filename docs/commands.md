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
  * `!addmodcommand !<name> <optional: -cd=seconds> <response>` - Adds a text command only usable by mods. (overwrites existing)
    * *-cd= sets a cooldown in seconds on the command*
    * *Example:* `!addcommand !usefulmessage Reminder that subscribing is not necessary to support the streamer!`
    * *Example:* `!addcommand !superspam -cd=30 SPAM SPAM Kappa SPAM`
  
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