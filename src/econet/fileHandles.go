package econet

/*
Only 255 file handles a allowed per session, as file handle is identified by a single byte on the client
File servers often only allowed 255 file handles total for the server, on this system we have 255 handles
per user session e.g. logged on user at a specific machine.

 * can support as many clients as we want
 * @param \HomeLan\FileStore\Authentication\User $oUser
 * @return int
*/

// Get the next free id for a file handle for the given user
