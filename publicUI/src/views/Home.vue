<template>

  <h1 class="p-text-secondary">Hablutzel Consulting provides a number of services</h1>

  <div class="panel-container">
    <panel header="Architecture Academy" @click="this.$router.push( {name: 'AcademyOverview'})">
      The Architecture Academy provides a structured approach to learning about architecture.
      No matter if you're a long time practitioner or just starting out in the architecture profession,
      there will be something in the Academy that you can learn from or leverage in your work.
    </panel>
    <divider layout="vertical"/>
    <panel header="Individual & Team Training">
      Sometimes you need a more interactive approach to learning about architecture. In those
      cases, it's valuable to have the help of a practitioner with years of experience in a variety
      of companies. This interactive training is highly customized to your needs, and is intended to
      teach you how to <i>think like an architect</i> by learning from one of the best.
    </panel>
    <divider layout="vertical"/>
    <panel header="Consulting Services">
      There are times when you simply need to have some help in managing your architecture - particularly
      when setting up a new team or looking to fix a struggling one. Our consulting engagements are
      focused on creating a high performing architecture team - something we've done repeatedly at multiple
      large companies.
    </panel>
  </div>



</template>

<script>

export default {
  name: "HomeComponent",
}
</script>

<style scoped>

</style>