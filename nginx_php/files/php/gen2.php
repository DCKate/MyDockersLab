<?php
    require_once __DIR__ . '/pmodules/CacheModule.php';
    if (isset($_SERVER['TOKEN']))
    {
        $tok = $_SERVER['TOKEN'];

        $red_ser=$_SERVER['REDIS_SERVER'];
        $red_port=$_SERVER['REDIS_PORT'];
        list($ret,$redis)=RedisConnect($red_ser, $red_port);
        if($ret==0){
            list($ret,$json)=RedisGetJson($redis,$tok);
            if($ret==0){
                $obj = json_decode($json);
                if (file_exists($obj->{'path'}))
                {
                    error_log("file_exist");
                    readfile($obj->{'path'});
                    return;
                }
            }
        }
        echo "Error !! file not found!";
        
    }
?>