<?php

if($_SERVER['REQUEST_METHOD'] == 'GET') :




else :

require_once('load.php');
require_once('admin-hooks-post.php');

function  cmt_header_json() {
	header('HTTP/1.1 200 Not OK');
	header('Content-Type: application/json');
}

function error($msg) {
	cmt_header_json();
	echo json_encode([
		'errno' => 'error',
		'error' => $msg
		]);
	die(-1);
}

function cmt_make_public(&$cmts) {
	global $tbopt;
	$admin_email = $tbopt->get('email');

	$flts = ['email', 'ip', 'agent', 'status'];
	for($i=0,$c=count($cmts); $i < $c; $i++) {
		$cmt = $cmts[$i];
		$cmt->avatar = md5(strtolower($cmt->email));
		$cmt->is_admin = strcasecmp($cmt->email, $admin_email)==0;
		foreach($flts as &$f) unset($cmt->$f);

		if(isset($cmt->children)) {
			for($x=0,$xc=count($cmt->children); $x < $xc; $x++) {
				$child = $cmt->children[$x];
				$child->avatar = md5(strtolower($child->email));
				$child->is_admin = strcasecmp($child->email, $admin_email)==0;
				foreach($flts as &$f) unset($child->$f);
			}
		}
	}
	return $cmts;
}

function cmt_get_cmt() {
	global $tbcmts;

	cmt_header_json();

	$cmts = cmt_make_public($tbcmts->get($_POST));
	
	echo json_encode([
		'errno'		=> 'success',
		'cmts'		=> $cmts,
		]);
	die(0);
}

function cmt_post_cmt() {
	global $tbcmts;
	global $tbdb;

	$ret_cmt = (int)(isset($_POST['return_cmt']) ? $_POST['return_cmt'] : '');
	
	$r = $tbcmts->insert($_POST);
	if(!$r) {
		header('HTTP/1.1 200 Not OK');
		header('Content-Type: application/json');

		echo json_encode([
			'errno'	=> 'error',
			'error' => $tbcmts->error
			]);
		die(-1);
	}

	// 如果评论成功，则保存用户的昵称和邮箱
	// 使用localStorage应该会更好，但可能还未广泛应用
	setcookie('tb_cmt_user',	$_POST['author'],	strtotime('+1 year'), '/');
	setcookie('tb_cmt_email',	$_POST['email'],	strtotime('+1 year'), '/');
	setcookie('tb_cmt_url',		$_POST['url'],		strtotime('+1 year'), '/');

	header('HTTP/1.1 200 OK');
	header('Content-Type: application/json');

	ob_start();
	if($ret_cmt) {
		$c = ['id'=>$r];
		$cmts = cmt_make_public($tbcmts->get($c));

		echo json_encode([
			'errno'	=> 'success',
			'cmt'	=> $cmts[0],
			]);
	} else {
		echo json_encode([
			'errno'	=> 'success',
			'id'	=> $r,
			]);
	}
	header('Content-Length: '.ob_get_length());
	header('Connection: close');
	ob_end_flush();
	fastcgi_finish_request();

	apply_hooks('comment_posted', 0, $_POST);
	die(0);
}

function cmt_get_count() {
	global $tbcmts;

	echo $tbcmts->get_count((int)$_POST['post_id']);

	die(0);
}

$do = isset($_POST['do']) ? $_POST['do'] : '';

if($do == 'get-cmt') cmt_get_cmt();
if($do == 'post-cmt') cmt_post_cmt();
if($do == 'get-count') cmt_get_count();

endif;

