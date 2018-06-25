<?php
    header('Content-Type: application/json');
    require_once __DIR__ . '/pmodules/CacheModule.php';
    date_default_timezone_set("UTC"); 

	session_start();
    $sessionID=session_id();
    error_log('ask sess'.$sessionID);
    $red_ser=$_SERVER['REDIS_SERVER'];
    $red_port=$_SERVER['REDIS_PORT'];

    define("FILE_QUERY_PATH", "/home/ubuntu/www/php/sess/");
    $file_name = "ask.txt";
    $query_folder = FILE_QUERY_PATH.$sessionID."/";
    $m3u8file = '/home/ubuntu/www/php/protected/demo/bwf.m3u8'; 
    if (!file_exists($query_folder))//check query session folder exist
    {
        mkdir($query_folder, 0777, true);
    }
    $motimestr = sprintf("%d",time());
    $moticket = hash('ripemd256',$motimestr.$sessionID.'mo.png');

    $file_path = FILE_QUERY_PATH.$sessionID."/".$file_name;
    $file = fopen($file_path, "w");
    //I suspect you did not know that there are different & escapes in HTML. 
    // The W3C you can see the codes. &times means × in HTML code. Use &amp;times instead.
    fwrite($file, "http://localhost:55555/seek.php?ticket=".$moticket."&amp;timestamp=".$motimestr);
    fclose($file);
    // $date = new DateTime('now');
    // $date->format('c');
    $rkey = hash('gost',$motimestr.$sessionID);
    $hlskey = hash('gost',$motimestr.$sessionID.'hlsrequest');
    $a = array('path' => $file_path, 'expired' => 180);
    $c = array('path' => $m3u8file, 'expired' => 300);
    list($ret,$redis)=RedisConnect($red_ser, $red_port);
    $rarry=array();
    if($ret==0){
        $ret=RedisSaveJson($redis,$rkey,json_encode($a));
        array_push($rarry,array('ret' => $ret, 'url' => urlencode("http://localhost:55555/gen.php?token=".$rkey)));
        $ret=RedisSaveJson($redis,$hlskey,json_encode($c));
        array_push($rarry,array('ret' => $ret, 'url' => urlencode("http://localhost:55555/".$hlskey.".m3u8")));
        echo urldecode(json_encode($rarry));
    }
?>