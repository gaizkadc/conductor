//
// Copyright (C) 2017 Daisho Group - All Rights Reserved
//
// List of messages - DEPRECATED.

package derrors

// DO NOT include constant messages in this file.
// In the initial version of the error repository error messages where stored in this file, with the idea that
// it will form a central repository for all error messages. However, this limits the addition of new messages
// as it will be necessary to update the derror repository and also may tempt the developer to use less significant
// messages to avoid creating new ones.
// Based on this, each repository will now define its own errors. Add a new directory called errors, with a file
// called messages.go to define the repository errors. Use always this layout so we can employ tools to adapt/translate
// those messages when needed.
