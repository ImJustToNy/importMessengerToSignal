# importMessengerToSignal

This project allows you to import your Facebook Messenger conversation into Signal. It relies on the [signal-cli](https://github.com/AsamK/signal-cli) project. Consider this project as an alpha version, it may not work as expected and may require some manual intervention.

## Usage
1. Obtain the copy of your Facebook Messenger data from Facebook by following the instructions [here](https://www.facebook.com/help/1701730696756992) and extract the contents of the downloaded archive
2. Clone this repository
3. Build the project with `go build .`
4. [Obtain the `signal-cli` binary for your system](https://github.com/AsamK/signal-cli/wiki/Binary-distributions)
5. Register all the phone numbers that take part in the conversations you want to import with `signal-cli` using `./signal-cli link`
6. Copy the example configuration file `config.example.yml` to `config.yml` and fill in the required fields
7. Run the program with `./importMessengerToSignal`