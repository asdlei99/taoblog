<?php

function tb_head() {
?>
<style>
    .left {
        float: left;
    }
    .right {
        float: right;
    }
    .list {
        list-style: none;
    }
    .item {
        margin-bottom: 3em;
        /*border: 1px solid gray;*/
        /*border-radius: 4px;*/
    }
    .meta {
        padding-bottom: 0.5em;
        border-bottom: 1px solid #e1e1e1;
    }
    .meta, .content {
        margin: 1em;
        overflow: hidden;
    }
</style>
<?php
}

add_hook('tb_head', 'tb_head');
require('header.php');

$posts = $tbpost->query_by_date(0,0,5);

echo "<ul class=list>\n";

foreach($posts as $post) {
?>
<li class=item>
    <div>
        <div class=meta>
            <div class=left>
            <?php
                echo htmlspecialchars($post->title);
            ?>
            </div>
            <div class=right>
            <?php
                echo $post->date;
             ?>
            </div>
        </div>
        <div class=content>
            <?php echo $post->content; ?>
        </div>
    </div>
</li>
<?php
}

echo "</ul>\n";

require('footer.php');
