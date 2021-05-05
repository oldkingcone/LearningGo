package bannerLib

import (
    `github.com/gookit/color`
)

func PrintMainBanner() bool{

    banner := `

▄▄▄▄    █    ██   ██████ ▄▄▄█████▓▓█████  ██▀███
▓█████▄  ██  ▓██▒▒██    ▒ ▓  ██▒ ▓▒▓█   ▀ ▓██ ▒ ██▒
▒██▒ ▄██▓██  ▒██░░ ▓██▄   ▒ ▓██░ ▒░▒███   ▓██ ░▄█ ▒
▒██░█▀  ▓▓█  ░██░  ▒   ██▒░ ▓██▓ ░ ▒▓█  ▄ ▒██▀▀█▄
░▓█  ▀█▓▒▒█████▓ ▒██████▒▒  ▒██▒ ░ ░▒████▒░██▓ ▒██▒
░▒▓███▀▒░▒▓▒ ▒ ▒ ▒ ▒▓▒ ▒ ░  ▒ ░░   ░░ ▒░ ░░ ▒▓ ░▒▓░
▒░▒   ░ ░░▒░ ░ ░ ░ ░▒  ░ ░    ░     ░ ░  ░  ░▒ ░ ▒░
░    ░  ░░░ ░ ░ ░  ░  ░    ░         ░     ░░   ░
░         ░           ░              ░  ░   ░
░



built by: oldkingcone
**********************************************
* Dont break the law.                        *
*                                            *
*                                            *
*                                            *
*                                            *
*                                            *
**********************************************
dont judge me, I'm learning go.
`
    color.Style{color.BgRed, color.FgWhite, color.Bold}.Printf(
        "%s",
        banner)
    return true
}