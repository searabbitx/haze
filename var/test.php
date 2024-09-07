<?php

if (isset($_GET['foo']) && $_GET['foo'] == "bar'") {
    http_response_code(500);
    echo "Error!";
} else {
    echo "Ok!";
}
