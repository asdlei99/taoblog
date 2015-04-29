<?php

require('header.php');
?>
<div class="query category-query">
	<ul class="item-list">
<?php
while($tbquery->has()){
	$the = $tbquery->the();
?>
	<li class="item cat-item"><h2><a target="_blank" href="<?php 
			if($the->type === 'post') {
				echo the_post_link($the, false);
			} else if($the->type === 'page') {
				echo the_page_link($the, false);
			}
			?>"><?php echo $the->title;?></a></h2></li>
<?php
} ?>
	</ul>
</div><?php

require('footer.php');

