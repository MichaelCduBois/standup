# Standup

Standup: Streamline daily stand-ups with a CLI tool that summarizes your previous day's work for quick updates. Boost team communication effortlessly!

## Usage Examples
```bash
# Output Standup Notes
standup

# Add item to standup notes
standup --add "Added a sweet feature."

# Add blocker to standup notes
standup --add --blocker "Waiting for code review."

# Add item to standup notes for previous day
standup --add --yesterday "Item I forgot to add."

# Age standup notes
standup --age

# List all standup notes
standup --list

# Delete item from standup notes
standup --delete "item key"

# Delete all items form standup notes
standup --reset
```
