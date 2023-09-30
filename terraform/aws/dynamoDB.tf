resource "aws_dynamodb_table" "this" {
	name           = "TicTactics"
	billing_mode   = "PAY_PER_REQUEST"
	hash_key       = "PK"

	attribute {
		name = "PK"
		type = "S"
	}

	attribute {
		name = "SK"
		type = "S"
	}

	attribute {
		name = "GSI1PK"
		type = "S"
	}

	attribute {
		name = "GSI1SK"
		type = "S"
	}

	global_secondary_index {
        name               = "GSI1"
        hash_key           = "GSI1PK"
        range_key          = "GSI1SK"
        projection_type    = "ALL"
        write_capacity     = 1
        read_capacity      = 1
    }
}
