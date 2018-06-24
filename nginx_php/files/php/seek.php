<?php
    session_start();
    $sessionID=session_id();
    error_log('seek sess'.$sessionID);
    if (isset($_GET['ticket'])&&isset($_GET['timestamp']))
    {
        $tic = $_GET["ticket"];
        $timestamp = $_GET["timestamp"];
        if($tic==hash('ripemd256',$timestamp.$sessionID.'mo.png')){
            $aliasedFile = '/download/mo.png'; //this is the nginx alias of the file path
            $realFile = '/home/ubuntu/www/php/protected/mo.png'; //this is the physical file path
            $filename = 'download.png'; //this is the file name user will get
            header('Cache-Control: public, must-revalidate');
            header('Pragma: no-cache');
            header('Content-Type: application\pdf');
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