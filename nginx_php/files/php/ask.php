<?php
    require_once __DIR__ . '/pmodules/CacheModule.php';
    date_default_timezone_set("UTC"); 

	session_start();
    $sessionID =session_id();
    
    $red_ser=$_SERVER['REDIS_SERVER'];
    $red_port=$_SERVER['REDIS_PORT'];

    define("FILE_QUERY_PATH", "/home/ubuntu/www/php/sess/");
    $file_name = "ask.txt";
    $query_folder = FILE_QUERY_PATH.$sessionID."/";
    if (!file_exists($query_folder))//check query session folder exist
    {
        mkdir($query_folder, 0777, true);
    }
    $file_path = FILE_QUERY_PATH.$sessionID."/".$file_name;
    $file = fopen($file_path, "w");
    fwrite($file, 'Hello Hello'.PHP_EOL.'You are goodgood');
    fclose($file);
    // $date = new DateTime('now');
    // $date->format('c');
    $rkey = hash('gost',sprintf("%d",time()).$sessionID);
    $a = array('path' => $file_path, 'expired' => 180);
    list($ret,$redis)=RedisConnect($red_ser, $red_port);
    if($ret==0){
        $ret=RedisSaveJson($redis,$rkey,json_encode($a));
        $b =array('ret' => $ret, 'url' => 'http://localhost:55555/gen.php?token='.$rkey);
        echo urldecode(json_encode($b));
    }

?>