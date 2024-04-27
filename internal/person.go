package internal

func getPersonByFacebookID(facebookID string) *Person {
	for _, person := range GetConfig().People {
		if person.FacebookID == facebookID {
			return &person
		}
	}

	return nil
}
