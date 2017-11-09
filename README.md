# kube-bot
A Kubernetes Slack bot. Addressed by the named `@kube`.  

# Table Of Contents
- [Commands](#commands)
- [Contributing](#contributing)

# Commands
You communicate with Kube in a command line-ish way. The first word of every 
message must be a valid command. These commands can have sub commands and 
arguments.  

Most commands have a `<query>` argument. Which specifies which Kubernetes 
resource you are trying to perform an action on. This query format follows the 
exact same syntax as the kubectl tool: `<type>/<name>[/<revision>]`.

Kube responds to the following commands:

✅ = Available, ⭕ = In Progress

- ⭕ `get <query>`
    - Retrieves a list of resources based on the given query
- ⭕ `describe <query`>
    - Shows detailed information about a particular item specified by the query
- ⭕ `rollout <sub command> <query>`
    - Performs actions related to application rollouts (updating versions, 
      deployments)
    - The application to perform rollout related actions on is specified by the 
      query option
    - Available sub commands are:
        - `status`
            - Shows the status of the specified rollout
        - `pause`
            - Pauses the specified rollout
        - `resume`
            - Resumes the specified rollout
        - `history`
            - Shows the rollout history for the specified resource
        - `undo`
            - Revert to the specified rollout
- ⭕ `scale <query> <number>`
    - Scales the resource specify by the provided query to the provided number 
      of replicas
- ⭕ `logs <query> [top | bottom] [<lines>]`
    - Displays the logs from the resource specified by the query
    - By default 25 lines from the logs are shown. This can be changed by 
      providing your own number of lines as an argument
    - By default this command will display the # of lines specified from the 
      bottom of the logs. This can be changed by providing the `top` or `bottom` 
      keywords (Not both)
- ⭕ `version`
    - Displays version information about kube-bot

All of which are existing kubectl commands. The following original commands 
are also provided:

- ⭕ `track <repo name> [#<channels...>] [anyones | mine | none] [all | failure | success] [<branches...>]`
    - Signs you up to receive messages about application deployments
    - The repo name argument should be the name of a GitHub repo. In the form 
      `username/repository`.
    - An optional list of channels to update about deployments can be provided. 
      Channels names must have a `#` before their name, and be separated by 
      spaces
    - One of, and no more, of the following keywords can be provided after the 
      repo name argument: `anyones`, `mine`, or `none`
        - If none of these keywords is provided `mine` is assumed
        - `anyones` signs you up for updates about all commits
        - `mine` signs you up for updates only about your commits
        - `none` unsubscribes you from any updates
    - One of the following keywords can be provided after the repo name
      argument:
        - If none of these keywords is provided `all` is assumed
        - `all` signs you up for updates about all deployments
        - `failure` signs you up to be updated only when a deployment fails
        - `success` signs you up to be updated only when a deployment succeeds 
    - The branches argument can be used to customize which branches you receive 
      updates about
        - If none are provided, you are signed up for updates about all branches 
        - Multiple branch names can be provided, separated by spaces
- ⭕ `config [<key>] [<value>]`
    - Manages kube-bot configuration
    - If no key is provided, a general settings message will be sent in 
      response. Where users can edit settings via a in Slack UI of sorts
    - One can view a key's value by providing only the key argument
    - If only the key argument is provided, but no parameter matches exactly, 
      a search will be performed, and similar settings keys will be listed
    - If the key and value arguments are provided, the specified option will be 
      set to the specified value
    - Valid keys are:
        - To be determined
- ⭕ `help`
    - Prints this command list

# Contributing
Kube Bot is an Open Source Software project. And welcomes contributions from 
anyone and everyone.  

Check out the [issue tracker](https://github.com/Noah-Huppert/kube-bot) to find 
work. Or open an issue if you found a bug or had an idea for a new feature. 
Or, even better, open a Pull Request with the code for that new feature!
