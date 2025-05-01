# Example: Create File Restore

## Summary
This example performs the following steps using the slide API:

- Checks for an API Token variable, and if not found, prompts the user for a token
- Queries the Slide API for a list of Agents
- Presents the list of Agents to the user, and requests the index of the agent that
	the restore should be created for
- Queries the Slide API for a list of snapshots for the provided Agent
- Creates a File Restore for the given Agent and Snapshot
- Provides the download URL for the File Restore.

## Example
![Create File Restore Example](../../ops/docs/assets/create_file_restore.gif)