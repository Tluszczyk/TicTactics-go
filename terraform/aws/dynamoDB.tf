resource "aws_dynamodb_table" "this" {
	name           = "example_table"
	billing_mode   = "PAY_PER_REQUEST"
	hash_key       = "id"
	stream_enabled = true

	attribute {
		name = "id"
		type = "S"
	}

	attribute {
		name = "username"
		type = "S"
	}

	attribute {
		name = "game_id"
		type = "S"
	}

	attribute {
		name = "password"
		type = "S"
	}

	global_secondary_index {
		name               = "username-index"
		hash_key           = "username"
		projection_type    = "ALL"
		write_capacity     = 1
		read_capacity      = 1
	}
}
