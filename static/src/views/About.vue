<template>
  <div style="display: flex; flex-direction: column; padding: 5px 5px 5px 5px;">
    <panel header="About Bob">
      <div style="display: flex">
        <img class="headshot" src='/images/headshot.JPG' alt="A headshot of Bob"/>
        <divider layout="vertical"/>
        <div style="display: flex; flex-direction: column;">
          <p>
            Bob Hablutzel has decades of architecture experience, including being leading the enterprise
            architecture practice at multiple Fortune 100 companies. His architectural approach is rooted
            in the practical, levering both his engineering and business background. He is available
            for consulting, training, and mentoring engagements. In addition, he has created the
            <router-link class="p-link p-text-secondary" :to="{name: 'AcademyOverview'}">Architecture Academy</router-link>
            to share his experiences with other
            architects and aspiring architects worldwide.
          </p>
          <a href="mailto:info@hablutzel.com">Contact via email</a>
        </div>
      </div>
    </panel>
    <panel header="About this site" style="margin-top: 5px">
      <p>
        The site is written in <a href="https://vuejs.org/" target=”_blank”>Vue3</a>
        and hosted on <a href="https://render.com/" target=”_blank”>render.com</a>. The site
        is built using an open-source content management system
        <a href="https://github.com/bobhablutzel/jutzo" target="_blank">Jutzo</a>
        if you are interested in learning how it was put
        together. [NOTE: Work in progress!]
      </p>
    </panel>
  </div>

</template>


<script>

export default {
  name: "AboutComponent",
}
</script>

<style scoped>

.headshot {
  width: 120px;
  height: 160px;
}


</style>