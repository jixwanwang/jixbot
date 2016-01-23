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
* add (*M*)
  * !addcommand !<name> <response> - Adds a text command (overwrites existing).
  * !addmodcommand !<name> <response> - Adds a text command only usable by mods (overwrites existing).
* delete (*M*)
  * !deletecommand !<name> - Deletes a text command.
* modonly (*M*)
  * !modonly - Makes all jixbot commands only usable by mods.
  * !subonly - Makes all jixbot commands only usable by subs.
  * !commandmode - Displays the restrictions on jixbot commands.

Stats
---------------
* uptime (*V*)
  * !uptime - Displays how long the stream has been online.
* lines_typed (*V*)
  * !linestyped - Whispers the caller how many lines they've typed in chat.
* time_spent (*V*)
  * !timespent - Whispers the caller how long they've been spent in online chat.
  * !longestwatcher - Displays the users with the highest viewing times.
* money (*V*)
  * !cash - Whispers the caller how much currency they have.
  * !givecash <username> <amount> - Gives <amount> coins to <username> .
  
Passive commands
----------------
* failfish - When someone mis-capitalizes the FailFish emote, corrects them
* conversation - When someone mentions "jixbot", jixbot replies
* combo - Combo amount increases as a particular emote or phrase is said in chat, kept alive by unique chatters who say it. Ends when enough time has passed since the last time combo was kept alive.
* submessage - Welcome message is displayed when someone subs. Uses sub emotes, so having jixbot subbed is useful.

Misc commands
---------------
* emotes (*V*)
  * !emotes - Lists emotes in the channel.
* brawl
  * !brawl (*M*) - Starts a brawl. Anyone can join with !pileon, only one person wins.
  * !pileon <optional: weapon> (*V*) - Join an active brawl with an optional weapon.
  * !brawlstats <season or "all"> (*V*) - Display top brawl winners of season or overall.
  * !newbrawlseason (*B*) - Start a new brawl season.