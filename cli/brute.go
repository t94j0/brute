package cli

import brute "github.com/t94j0/brute/brute-protocol"

func createWordlist(wordlist *WordListBrute, input chan<- brute.BruteData) {
	nextHost := wordlist.NextHost()
	for nextHost != "" {
		hasUsername, nextUsername := wordlist.NextUsername()
		for hasUsername {
			hasPassword, nextPassword := wordlist.NextPassword()
			for hasPassword {
				bruteChan <- brute.BruteData{
					nextHost,
					nextUsername,
					nextPassword,
				}
				hasPassword, nextPassword = wordlist.NextPassword()
			}
			wordlist.PasswordList.Reset(wordlist.PasswordInitial)
			hasUsername, nextUsername = wordlist.NextUsername()
		}
		wordlist.UsernameList.Reset(wordlist.UsernameInitial)

		nextHost = wordlist.NextHost()
	}
}

func Brute(wordlist *WordListBrute, module brute.Brute) error {
	bruteChan := make(chan<- brute.BruteData)
	resultsChan := make(<-chan brute.BruteData)
	go createWordlist(wordist, bruteChan)
	brute.BruteChan(module, bruteChan, results)
}
