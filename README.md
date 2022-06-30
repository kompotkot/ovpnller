# ovpnller

## Example of default configuration

```json
{
	"ca": {
		"address": "127.0.0.1",
		"port": "22",
		"username": "ubuntu",
		"private_key_path": "/home/ubuntu/.ssh/id_rsa"
	},
	"server": {
		"address": "127.0.0.1",
		"port": "22",
		"username": "ubuntu",
		"private_key_path": "/home/ubuntu/.ssh/id_rsa"
	},
	"actions": {
		"ca_init": [
			{
				"action": "EASYRSA_PKI='/home/ubuntu/easy-rsa/pki' bash /home/ubuntu/easy-rsa/easyrsa build-ca nopass",
				"action_type": "command"
			}
		],
		"client_register": [],
		"server_register": [
			{
				"action": "bash /home/ubuntu/test.bash",
				"action_type": "command",
				"source_file_path": "",
				"target_file_path": ""
			}
		]
	}
}
```
