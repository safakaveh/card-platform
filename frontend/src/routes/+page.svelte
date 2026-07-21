<script lang="ts">
	import * as Card from "$lib/components/ui/card/index.js";
</script>

<Card.Root class="mx-auto w-1/2">
 <Card.Header>
  <Card.Title>به برنامه دستگاه XYZ خوش آمدید</Card.Title>
  <Card.Description>پیاده سازی و طراحی شرکت تراشه هوشمند</Card.Description>
 </Card.Header>
 <Card.Content class="-mb-(--card-spacing)">
  <p class="mb-2">در هنگام استفاده از برنامه به موارد زیر توجه فرمایید</p>
  
  <ol>
    <li>فایل بارگذاری حتما CSV استاندارد با header باشد</li>
    <li>در صورتی که فایل excle دارید کافی است آن را به فرمت .csv ذخیره کنید</li>
    <li>در صورتی که header با حروف frn_ شروع شده باشد موارد آن ستون در روی کارت لیزر خواهند شد</li>
    <li>در صورتی که header با حروف bck_ شروع شده باشد موارد آن ستون در پشت کارت لیزر خواهند شد</li>
    <li>در صورتی که قصد لیزر UID کارت ها را دارید کافیست یک ستون خالی با  هدر frn_uid یا bck_uid داشته باشید </li>
  </ol>
 </Card.Content>
 <Card.Footer class="mt-7 text-center justify-center">
  <p>در آخر از توجه شما سپاسگزاریم و امیدواریم از استفاده از دستگاه ما رضایت کامل داشته باشید</p>
 </Card.Footer>
</Card.Root>



