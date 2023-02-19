# Cinnamon

Cinnamon is a discord bot in development, planned to have many features, as well as focus on self-hosting, while also providing a public hosted bot for those who don't want to deal with hosting fees.
Some of its core features will include:
- Golang
- SQLite use possible, however, comes with heavy performance penalties, particularly with database heavy workloads. Recommended only for single server use cases. 
- Golang
- Flagship feature being designed with an encouragement for self hosting and customizimg one's own instance.
 - No need to host your own dashboard, our own dashboard will allow you to generate a dashboard API key with which you can link your own instance to our dashboard. 
 - Dashboard design is not hard coded, each bot module can contain a dashboard definition file, allowing it to create its own configuration tab, containing multiple sub-tabs and allowing you to define what each page should contain, such as text boxes, drop down lists, multiple choice boxes, ect. Not unlike discord components. 
 - You can choose which modules you want to compile with the bot, even allowing to omit all modules, including the dashboard connection module. However, there's still a few core features that cannot be removed, being the websocket handler, and a few core commands. 
 - You can even use the included websocket handler for your own modules, allowing you to add your own handlers.
 - A flagship feature will be the ability for a user to migrate one's own data from the main Cinnamon instance to other 3rd party Cinnamon-based instances. This will require the user to verify via OAuth and consent to transfer their information from the base modules included with the bot. This might even allow for an instance with a custom module to advertise itself to other instances with the same custom module as able to transfer information between 3rd party instances. All this will be optional to disable. 

- A focus on non predatory monetization, with non essential, but compute heavy features locked behind pay wall. However, self hosted instances do not have this limit. 
- Placeholder
- ...
- Did I mention Golang?


Development of this bot is expected to take a while, being mostly started as a passion project, now being developed by two developers in their free time. 
