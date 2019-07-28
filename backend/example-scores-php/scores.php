<?php

// This file is an example of a backend written in PHP.

$handlers = [
	'add ? to ?' => function($args) use (&$total, &$count) {
        $total += $args[0];
        ++$count;
	},
	'average of ? into ?' => function($args) use (&$total, &$count) {
		return ["set" => ['$1' => (string)($total / $count)]];
	},
	'display ?' => function() use (&$total) {
		return ["text" => "The total is $total."];
	}
];

// The code following should not need to be changed.

$socket = socket_create(AF_INET, SOCK_STREAM, 0);
$result = socket_bind($socket, "127.0.0.1", $_ENV['BENTO_PORT']);
$result = socket_listen($socket, 3);
$spawn = socket_accept($socket);

while ($message = json_decode(socket_read($spawn, 65536, PHP_NORMAL_READ))) {
    if ($message->special === "sentences") {
    	$result = ['sentences' => array_keys($handlers)];
    } else {
    	$handler = $handlers[$message->sentence];
    	$result = $handler($message->args);
    }

    if (!$result) {
    	$result = new stdClass();
    }

    $output = json_encode($result) . "\n";
    socket_write($spawn, $output, strlen($output));
}

socket_close($spawn);
socket_close($socket);
