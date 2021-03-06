<?php

require_once('admin.php');

function postmanage_admin_head() {
?>

<style>

table, td, th {
    border: 1px solid gray;
    border-collapse: collapse;
}

th, td {
    padding: 8px;
}

td.title {
    max-width: 300px;
}

</style>

<script src="https://cdn.jsdelivr.net/npm/vue"></script>

<?php }

add_hook('admin_head', 'postmanage_admin_head');

admin_header();

?>

<table id="table">
<thead>
<tr>
<th>编号</th>
<th>标题</th>
<th>发表日期</th>
<th>修改日期</th>
<th>浏览量</th>
<th>评论数</th>
<th>源类型</th>
<th>操作</th>
</tr>
</thead>
<tbody>
  <tr v-for="p in posts">
    <td>{{ p.id }}</td>
    <td>{{ p.title }}</td>
    <td>{{ p.date }}</td>
    <td>{{ p.modified }}</td>
    <td>{{ p.page_view }}</td>
    <td>{{ p.comment_count }}</td>
    <td>{{ p.source_type }}</td>
    <td><a target=_blank v-bind:href="'/admin/post.php?do=edit&amp;id=' + p.id">编辑</a></td>
  </tr>
</tbody>
</table>

<script>
$.get('/v1/posts!manage', function(posts) {
    new Vue({
        el: '#table',
        data: {
            posts: posts,
        },
    })
});
</script>

<?php
admin_footer();
