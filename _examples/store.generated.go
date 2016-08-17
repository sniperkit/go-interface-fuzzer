package example

import (
	"errors"
	"fmt"
	"model"
	"rand"
	"reflect"
	"testing"
)

// Store

func FuzzTestStore(makeTest func(int) Store, t *testing.T) {
	rand := rand.New(rand.NewSource(0))

	err := FuzzStore(makeTest, rand, 100)

	if err != nil {
		t.Error(err)
	}
}

func FuzzStore(makeTest func(int) Store, rand *rand.Rand, max uint) error {
	var (
		argInt int
	)

	argInt = rand.Int()

	expectedStore := makeReferenceStore(argInt)
	actualStore := makeTest(argInt)

	return FuzzStoreWith(&expectedStore, actualStore, rand, max)
}

func FuzzStoreWith(reference Store, test Store, rand *rand.Rand, maxops uint) error {
	// Create initial state
	state := uint(0)

	for i := uint(0); i < maxops; i++ {
		// Pick a random number between 0 and the number of methods of the interface. Then do that method on
		// both, check for discrepancy, and bail out on error. Simple!

		actionToPerform := rand.Intn(7)

		switch actionToPerform {
		case 0:
			// Call the method on both implementations
			var (
				argModelIDMessage model.IDMessage
			)

			argModelIDMessage, state = generateIDMessage(rand, state)

			expectedError := reference.Put(argModelIDMessage)
			actualError := test.Put(argModelIDMessage)

			// And check for discrepancies.
			if !((expectedError == nil) == (actualError == nil)) {
				return fmt.Errorf("inconsistent result in Put\nexpected: %v\nactual:   %v", expectedError, actualError)
			}
		case 1:
			// Call the method on both implementations
			var (
				argModelID      model.ID
				argModelChannel model.Channel
			)

			argModelID, state = generateID(rand, state)
			argModelChannel = generateChannel(rand)

			expectedModelID, expectedModelIDMessage := reference.EntriesSince(argModelID, argModelChannel)
			actualModelID, actualModelIDMessage := test.EntriesSince(argModelID, argModelChannel)

			// And check for discrepancies.
			if !reflect.DeepEqual(expectedModelID, actualModelID) {
				return fmt.Errorf("inconsistent result in EntriesSince\nexpected: %v\nactual:   %v", expectedModelID, actualModelID)
			}
			if !reflect.DeepEqual(expectedModelIDMessage, actualModelIDMessage) {
				return fmt.Errorf("inconsistent result in EntriesSince\nexpected: %v\nactual:   %v", expectedModelIDMessage, actualModelIDMessage)
			}
		case 2:
			// Call the method on both implementations
			var (
				argModelID      model.ID
				argModelChannel model.Channel
			)

			argModelID, state = generateID(rand, state)
			argModelChannel = generateChannel(rand)

			expectedModelID, expectedMessageIterator := reference.EntriesSinceIter(argModelID, argModelChannel)
			actualModelID, actualMessageIterator := test.EntriesSinceIter(argModelID, argModelChannel)

			// And check for discrepancies.
			if !reflect.DeepEqual(expectedModelID, actualModelID) {
				return fmt.Errorf("inconsistent result in EntriesSinceIter\nexpected: %v\nactual:   %v", expectedModelID, actualModelID)
			}
			if !compareMessageIterators(expectedMessageIterator, actualMessageIterator) {
				return fmt.Errorf("inconsistent result in EntriesSinceIter\nexpected: %v\nactual:   %v", expectedMessageIterator, actualMessageIterator)
			}
		case 3:
			// Call the method on both implementations
			expectedModelID := reference.MostRecentID()
			actualModelID := test.MostRecentID()

			// And check for discrepancies.
			if !reflect.DeepEqual(expectedModelID, actualModelID) {
				return fmt.Errorf("inconsistent result in MostRecentID\nexpected: %v\nactual:   %v", expectedModelID, actualModelID)
			}
		case 4:
			// Call the method on both implementations
			expectedInt := reference.NumEntries()
			actualInt := test.NumEntries()

			// And check for discrepancies.
			if !reflect.DeepEqual(expectedInt, actualInt) {
				return fmt.Errorf("inconsistent result in NumEntries\nexpected: %v\nactual:   %v", expectedInt, actualInt)
			}
		case 5:
			// Call the method on both implementations
			expectedModelIDMessage := reference.AsSlice()
			actualModelIDMessage := test.AsSlice()

			// And check for discrepancies.
			if !reflect.DeepEqual(expectedModelIDMessage, actualModelIDMessage) {
				return fmt.Errorf("inconsistent result in AsSlice\nexpected: %v\nactual:   %v", expectedModelIDMessage, actualModelIDMessage)
			}
		case 6:
			// Call the method on both implementations
			expectedInt := reference.MessageLimit()
			actualInt := test.MessageLimit()

			// And check for discrepancies.
			if !reflect.DeepEqual(expectedInt, actualInt) {
				return fmt.Errorf("inconsistent result in MessageLimit\nexpected: %v\nactual:   %v", expectedInt, actualInt)
			}
		}

		if !(reference.NumEntries() == len(reference.AsSlice())) {
			return errors.New("invariant violated: %var.NumEntries() == len(%var.AsSlice())")
		}

	}

	return nil
}
