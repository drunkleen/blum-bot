## Blum Bot
---

### Overview
This bot is designed to automate various tasks in a Telegram mini-app BlumBot. It handles everything from playing drop games, claiming referral balances, starting missions, and claiming mission balances once they are completed.

### Features
 - Automated Task Execution: Automatically handles all in-game tasks.
 - Drop Game Participation: Plays the drop games to earn rewards.
 - Referral Balance Claiming: Claims balances from referrals.
 - Mission Management: Starts missions and claims the mission balance once completed.

### Getting Started
#### Prerequisites
Before you begin, ensure you have met the following 
requirements:
- You have a Telegram account
- You have to download the executable of build it from source

How to get query? Follow this steps:
1. open Telegram Desktop and login.
2. go to `Settings`
3. click on `Advanced`
4. click on `Experimental settings`
5. turn on `Enable webview inspecting`
6. go your BlumBot page and click on `Launch Blum`
7. right click on the blum page and select `Inspect Element`
8. click on `Application` tab
9. in the left bar use the dropdown next to `session storage` and click on `telegram.blum.codes`
10. select `telegram_initparam`
11. right click on `tgwebappdata` and select `copy value`
12. create a new directory in the same directory as executable file and name it `configs`
13. create a new blank file with the name `query_list.conf`
14. paste your query IDs
15. run the bot and enjoy!


### Disclaimer

This bot is intended for educational purposes only. Use it at your own risk. The author is not responsible for any misuse or damage caused by using this bot.

### Contributing

Contributions are welcome! Please follow these steps to contribute:

1. Fork the repository.
2. Create a new branch (git checkout -b feature/your-feature).
3. Make your changes and commit them (git commit -m 'Add some feature').
4. Push to the branch (git push origin feature/your-feature).
    Create a pull request.

### License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
