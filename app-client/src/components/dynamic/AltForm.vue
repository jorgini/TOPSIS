<script setup>
  const alt = defineModel();
</script>

<script>
  export default {
    emits: ['delete-alt'],
    data() {
      return {
        isValidTitle: true,
        isVisitableDisc: false
      }
    },
    methods: {
      validate() {
        this.isValidTitle = this.modelValue.title.length > 0 && this.modelValue.title.length < 101;
      },
      deleteAlt() {
        this.$emit('delete-alt', null);
      },
      switchDisc() {
        this.isVisitableDisc = !this.isVisitableDisc
      }
    }
  }
</script>

<template>
  <div class="alt">
    <div class="row-cols-3">
      <div class="col-3">
        <p>Название:</p>
      </div>
      <div class="col-6">
        <input type="text" :class="{field: true, invalid: !isValidTitle}" name="title"
               placeholder="title" maxlength="100" v-model="alt.title" @input="validate" required/>
      </div>
      <div class="col-3 right-col">
        <img alt="" src="/cancel.png" class="cancel" @click="deleteAlt">
      </div>
    </div>
    <div class="row-cols-3">
      <div class="col-3">
        <img alt="" src="/arrow.png" :class="{flag: true, close: !isVisitableDisc}" @click="switchDisc">
        <p>Описание:</p>
      </div>
      <div class="col-6">
        <textarea type="text" :class="{field: true, invisible: !isVisitableDisc}" name="description"
            placeholder="description" maxlength="1000" v-model="alt.description"/>
      </div>
      <div class="col-3"></div>
    </div>
  </div>
</template>

<style scoped>
  @import "../../style.css";
  @import "../../assets/alts.css";
</style>