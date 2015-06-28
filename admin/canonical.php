<?php

function the_link(&$p, $home=true) {
	global $tbopt;
	global $tbtax;
	global $tbopt;

	$home = $home ? $tbopt->get('home') : '';
	$link = '';

	if($p->type === 'post') {
		$cats = implode('/', $tbtax->tree_from_id($p->taxonomy)['slug']);
		$slug = $p->slug;

		$link = $home.'/'.$cats.'/'.$slug.'.html';
	} else if($p->type === 'page') {
		$link = $home.'/'.$p->slug;
	} else {
		$link = '/';
	}

	return $link;
}

function the_edit_link(&$p, $ret_anchor = true) {
	$link = '/admin/post.php?do=edit&id='.(int)$p->id.'&type='.$p->type;

	return $ret_anchor
		? '<a href="'.$link.'">编辑</a>'
		: $link;
}

