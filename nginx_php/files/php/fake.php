<?php
    date_default_timezone_set("UTC"); 
    session_start();
    $sessionID=session_id();
    error_log('fake sess'.$sessionID);
    if (isset($_SERVER['bwf'])&&isset($_GET['ticket'])&&isset($_GET['timestamp']))
    {
        $bwf = $_SERVER["bwf"].'.ts';
        $tic = $_GET["ticket"];
        $timestamp = $_GET["timestamp"];
        if($tic==hash('crc32b',$timestamp.$sessionID.$bwf)){
            $aliasedFile = '/download/demo/'.$bwf; //this is the nginx alias of the file path
            $realFile = '/home/ubuntu/www/php/protected/demo/'.$bwf; //this is the physical file path
            $filename = $bwf; //this is the file name user will get
            header('Cache-Control: public, must-revalidate');
            header('Pragma: no-cache');
            header('Content-Type: video/MP2T');
            header('Content-Length: ' .(string)(filesize($realFile)) );
            header('Content-Disposition: attachment; filename='.$filename.'');
            header('Content-Transfer-Encoding: binary');
            header('X-Accel-Redirect: '. $aliasedFile);
            return;
        }
    }
    http_response_code(404);
    return;
?>